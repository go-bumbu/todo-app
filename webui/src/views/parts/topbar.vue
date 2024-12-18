<script setup>
import Horizontal from '@/lib/layout/horizontal.vue'
import UserProfile from '@/lib/user/UserProfile.vue'
import { useUserStore } from '@/lib/user/userstore.js'
import Button from 'primevue/button'

const user = useUserStore()
</script>

<template>
    <horizontal :center-content="true">
        <template v-slot:left>
            <router-link to="/app" class="layout-topbar-logo"> </router-link>
        </template>

        <!--        <InputText placeholder="Search" type="text" class="w-32 sm:w-auto" />-->
        <template v-slot:right>
            <UserProfile v-if="user.isLoggedIn" />
            <router-link v-if="!user.isLoggedIn" to="/login" class="layout-topbar-logo">
                <Button label="Login" icon="pi pi-sign-in" />
            </router-link>
        </template>
    </horizontal>
</template>
