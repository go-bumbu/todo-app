<script setup>
import Vertical from '@/lib/layout/Vertical.vue'
import Horizontal from '@/lib/layout/horizontal.vue'
import TopBar from '@/views/parts/topbar.vue'
</script>

<template>
    <div class="landing-page">
        <vertical :center-content="true">
            <template v-slot:header>
                <TopBar />
            </template>

            <template v-slot:main>
                <horizontal :center-content="true">
                    <template v-slot:default>
                        <div class="logo" />
                        <h1>Go Bumbu</h1>
                        <div>
                            This is a sample application that shows all of the features of Go Bumbu
                        </div>
                        <div>
                            Log in to the TODO app with <b>admin:admin</b> or <b>demo:demo</b>
                        </div>

                        <div>
                            ...or explore some other features on the
                            <a href="/demo"> demo page</a> (not served by the SPA)
                        </div>
                        <router-view />
                    </template>
                </horizontal>
            </template>
        </vertical>
    </div>
</template>
<style scoped>
.logo {
    width: 300px;
    height: 300px;
    background-image: url('/Logo.png');
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
}

.landing-page {
    background-color: #dfdcd7;
}
</style>
