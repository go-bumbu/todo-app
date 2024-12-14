<script setup>
import Avatar from 'primevue/avatar'
import Drawer from 'primevue/drawer'
import Button from 'primevue/button'
import { ref } from 'vue'
import { useUserStore } from '@/lib/user/userstore.js'
import { useTaskStore } from '@/stores/taskStore.js'
import router from "@/router/index.js";

const user = useUserStore()
const tasks = useTaskStore()
const visible = ref(false)

const logOut = () => {
    tasks.reset()
    user.logout(function (){
      router.push("/login")
    })
}
</script>
<template>
    <Avatar
        icon="pi pi-user"
        class="mr-2"
        size="large"
        @click="visible = true"
        style="background-color: #ece9fc; color: #2a1261; cursor: pointer"
    />

    <Drawer v-model:visible="visible" :header="user.loggedInUser" style="width: 25rem" position="right">
        <Avatar
            icon="pi pi-user"
            class="mr-3"
            size="xlarge"
            @click="visible = true"
            style="background-color: #ece9fc; color: #2a1261"
        />

        <p><Button label="Settings" icon="pi pi-cog" /></p>
        <p>
            <Button label="Logout" severity="danger" icon="pi pi-sign-out" @click="logOut()" />
        </p>
    </Drawer>
</template>
