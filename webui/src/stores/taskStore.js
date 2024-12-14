import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import axios from 'axios'

const taskEndpoint = import.meta.env.VITE_SERVER_URL_V0 + '/task'
const tasksEndpoint = import.meta.env.VITE_SERVER_URL_V0 + '/tasks'

export const useTaskStore = defineStore('tasks', () => {
    const tasks = ref([])
    const totalTasks = ref(0)
    const isDataLoaded = ref(false)
    const isLoading = ref(false)

    const processTasks = (payload) => {
        if (payload.Tasks && Array.isArray(payload.Tasks)) {
            // Iterate over each task in the payload
            payload.Tasks.forEach((task) => {
                insertTask(task.id, task.text, task.done)
            })
        } else {
            console.log('No tasks found in the payload.')
        }
    }
    const insertTask = (id, text, done) => {
        tasks.value.push({ task: text, id: id, done: done })
        totalTasks.value = tasks.value.length
    }
    const removeTask = (id) => {
        const index = tasks.value.findIndex(item => item.id === id)
        tasks.value.splice(index,1)
    }
    const updateTask = (id, text, done) =>{
        const index = tasks.value.findIndex(item => item.id === id)
        if (text !== ""){
            tasks.value[index].text = text
        }
        if (typeof done === "boolean"){
            tasks.value[index].done = done
        }
    }


    // loadTasks loads tasks from the remote api endpoint
    const loadTasks = () => {
        if (isDataLoaded.value === false) {
            isDataLoaded.value = true
            isLoading.value = true
            axios
                .get(tasksEndpoint)
                .then((res) => {
                    if (res.status === 200) {
                        processTasks(res.data)
                    } else {
                        console.log('err')
                        console.log(res)
                        // error?
                    }
                })
                .catch((err) => {
                    console.log(err)
                }).finally(()=>{
                    isLoading.value =false
            })
        }
    }
    // reset will wipe the store, this is used when a user logs out to prevent
    // data leak to another user on the same browser
    const reset = () => {
        tasks.value = []
        totalTasks.value = 0
        isDataLoaded.value = false
    }

    const addTask = (taskText) => {
        const taskPayload = {
            text: taskText
        }
        return axios
            .post(taskEndpoint, taskPayload, {
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then((res) => {
                console.log(res)
                if (res.status === 200) {
                    const payload = res.data
                    if (payload) {
                        insertTask(payload.id, payload.text, payload.done)
                    }
                } else {
                    console.log('err')
                    console.log(res)
                    // error?
                }
            })
            .catch((err) => {
                console.log(err)
            })
    }



    const deleteTask = (id) =>{
        return axios
            .delete(taskEndpoint+ "/"+id,)
            .then((res) => {
                if (res.status === 202) {
                    removeTask(id)
                } else {
                    console.log('err')
                    console.log(res)
                    // error?
                }
            })
            .catch((err) => {
                console.log(err)
            })
    }

    const setDone = (id, status) =>{
        const taskPayload = {
            id: id,
            done: status
        }
        return axios
            .put(taskEndpoint+ "/"+id, taskPayload, {
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then((res) => {
                if (res.status === 202) {
                    updateTask(id,"",status)
                } else {
                    console.log('err')
                    console.log(res)
                    // error?
                }
            })
            .catch((err) => {
                console.log(err)
            })
    }

    return {
        tasks, // list of tasks
        totalTasks, // amount of total tasks
        isLoading, // if the store is in a loading state
        loadTasks, // initial load of tasks
        addTask, // add a task
        setDone,
        deleteTask,
        reset // reset on logut
    }
})
