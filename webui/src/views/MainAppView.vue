<script setup>
import Horizontal from '@/lib/layout/horizontal.vue'
import Vertical from '@/lib/layout/Vertical.vue'
import TopBar from '@/views/parts/topbar.vue'
import { useUserStore } from '@/lib/user/userstore.js'
import TaskList from '@/views/parts/TaskList.vue'

const user = useUserStore()
</script>

<template>
    <vertical :center-content="false">
        <template v-slot:header>
            <TopBar />
        </template>

        <template v-slot:main>
            <horizontal :center-content="true">
                <template v-slot:default>
                    <TaskList />
                    <router-view />
                </template>
            </horizontal>
        </template>

        <template v-slot:footer>my footer </template>
    </vertical>
</template>
