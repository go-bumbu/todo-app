<script setup>
import Vertical from "@/lib/layout/Vertical.vue";
import Horizontal from "@/lib/layout/horizontal.vue";
import UserLogin from '@/lib/user/UserLogin.vue'
</script>

<template>
    <vertical :center-content="true">
        <template v-slot:main>
            <horizontal :center-content="true">
                <template v-slot:default>
                    <UserLogin />
                </template>
            </horizontal>
        </template>

        <template v-slot:footer> </template>
    </vertical>
</template>


