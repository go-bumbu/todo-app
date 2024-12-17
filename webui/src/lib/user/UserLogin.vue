<script setup>
import Card from 'primevue/card'
import Password from 'primevue/password'
import InputGroup from 'primevue/inputgroup'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import InputGroupAddon from 'primevue/inputgroupaddon'
import Button from 'primevue/button'
import Message from 'primevue/message'
import { ref } from 'vue'
import { useUserStore } from '@/lib/user/userstore.js'
import { Form } from '@primevue/forms';
import router from "@/router/index.js";
import LoadingScreen from "@/lib/loadingScreen.vue";
const user = useUserStore()

const userRef = ref(null)
const passRef = ref(null)
const keepLoginRef = ref(false)

const load = () => {
    user.login(
        userRef.value,
        passRef.value,
        keepLoginRef.value,
  function (){
      router.push('/app')
    })
}

</script>
<template>
    <Card>
        <template #title>Log in</template>
        <template #content>
          <Form class="">
          <div v-focustrap class="flex flex-column items-center gap-4">

              <InputGroup>
                    <InputGroupAddon>
                        <i class="pi pi-user"></i>
                    </InputGroupAddon>
                    <InputText placeholder="Username" v-on:keyup.enter="load" v-model="userRef" autocomplete="username" required />
                </InputGroup>

                <InputGroup>
                    <InputGroupAddon>
                        <i class="pi pi-lock"></i>
                    </InputGroupAddon>
                    <Password
                        v-model="passRef"
                        v-on:keyup.enter="load"
                        placeholder="Password"
                        :feedback="false"
                        toggleMask
                        :inputProps="{ autocomplete: 'current-password', required: true, }"
                    />
                </InputGroup>
              <InputGroup>
                <ToggleSwitch  inputId="loggedin" v-model="keepLoginRef" />
                <label  class="ml-4" for="loggedin">Keep me signed in</label>
             </InputGroup>
            <Message v-if="user.wrongPwErr" severity="error" closable
                    >Wrong username or password</Message
                >
                <Button label="Log in" class="w-full" @click="load" />
          </div>
          </Form>
        </template>
    </Card>
    <loadingScreen v-if="user.isLoading" />
</template>
