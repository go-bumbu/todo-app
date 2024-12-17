import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'


const statusPath = import.meta.env.VITE_AUTH_PATH + '/status'
const loginPath = import.meta.env.VITE_AUTH_PATH + '/login'
const logoutPath = import.meta.env.VITE_AUTH_PATH + '/logout'

export const useUserStore = defineStore('user', () => {

    const isLoggedIn = ref(false)
    const loggedInUser = ref("")
    const isFirstLogin = ref(true)
    const isLoading = ref(false)
    const wrongPwErr = ref(false)


    const user = ref("")

    const setFirstLoginFalse = () =>{
        isFirstLogin.value = false
    }

    const checkState = () =>{
        return axios
            .get(statusPath)
            .then((res) => {
                if (res.status === 200) {
                    console.debug(res.data)
                    loggedInUser.value = res.data['username']
                    isLoggedIn.value = res.data['logged-in']
                } else {
                    console.log(res)
                }
            })
            .catch((err) => {
                console.log(err)
            })
    }

    const login = (user, pass,keepMeLoggedIn, onSuccessNavigate) => {
        if (!keepMeLoggedIn){
            keepMeLoggedIn = false
        }
        const data = {
            username: user,
            password: pass,
            sessionRenew:keepMeLoggedIn,
        }

        const authAxios = axios.create()
        authAxios.interceptors.response.use(
            (response) => {
                return response
            },
            (error) => {
                if (error.response.status === 401) {
                    console.log('authentication returned 401')
                    isLoggedIn.value = false
                    wrongPwErr.value = true
                }
                return error
            }
        )
        isLoading.value = true

        authAxios
            .post(loginPath, data)
            .then((res) => {
                console.log(res)
                if (res.status === 200) {
                    console.debug(res.data)
                    loggedInUser.value = user
                    isLoggedIn.value = true
                    wrongPwErr.value = false

                    // Trigger the callback for navigation
                    if (onSuccessNavigate) {
                        onSuccessNavigate();
                    }

                } else {
                    console.log(res)
                }
            })
            .catch((err) => {
                console.log(err)
                // todo propagate login error

                // this.$toasted.show(
                //     'Please enter the correct details and try again',
                //     err,
                //     {
                //         position: 'top-left',
                //         duration: 200,
                //         type: danger,
                //     }
                // )
            })
            .finally(() => {
                isLoading.value = false
            })
    }

    const logout = (onLogout) =>{
        isLoading.value = true
        axios
            .post(logoutPath, '')
            .then((res) => {
                loggedInUser.value = ""
                isLoggedIn.value = false

                if (onLogout) {
                    onLogout();
                }

                // router.push('/login')
            })
            .catch((err) => {
                console.log(err)
                // todo propagate login error
            })
            .finally(() => {
                isLoading.value = false
            })
    }


    return {
        isFirstLogin,
        setFirstLoginFalse,

        isLoggedIn,
        loggedInUser,
        checkState,

        isLoading,
        login,
        wrongPwErr,

        logout

    }
})