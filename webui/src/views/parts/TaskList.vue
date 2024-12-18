<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Fluid from 'primevue/fluid'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Row from 'primevue/row'

import { ref, onMounted } from 'vue'
import { useTaskStore } from '@/stores/taskStore.js'
import LoadingScreen from '@/lib/loadingScreen.vue'

const tasks = useTaskStore()

const newTaskRef = ref('')
const addTask = () => {
    tasks.addTask(newTaskRef.value)
    newTaskRef.value = ''
}
const delTask = (id) => {
    console.log(id)
    // loading.value = true
    tasks.deleteTask(id)
    // setTimeout(() => {
    //     loading.value = false
    // }, 2000)
}
const toggleDone = (id, status) => {
    if (status === false) {
        tasks.setDone(id, true)
    } else {
        tasks.setDone(id, false)
    }
}

onMounted(() => {
    tasks.loadTasks()
})
const loading = ref(false)
const products = tasks.tasks
</script>

<style>
.taskDone {
    text-decoration: line-through;
}
</style>
<template>
    <Fluid class="inline-flex">
        <InputText
            id="username"
            v-on:keyup.enter="addTask"
            v-model="newTaskRef"
            placeholder="New task"
        />
        <Button icon="pi pi-send" @click="addTask" />
    </Fluid>

    <DataTable :value="products" tableStyle="width: 50rem; margin: 1rem">
        <Column headerStyle="display: none">
            <template #body="slotProps">
                <Button
                    icon="pi pi-check-circle"
                    text
                    aria-label="Filter"
                    @click="toggleDone(slotProps.data.id, slotProps.data.done)"
                />
            </template>
        </Column>

        <Column field="task" header="Task" headerStyle="display: none" #body="slotProps">
            <span :class="{ taskDone: slotProps.data.done === true }">
                {{ slotProps.data.task }}
            </span>
        </Column>
        <Column header="Delete" headerStyle="display: none">
            <template #body="slotProps">
                <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    aria-label="Cancel"
                    :loading="loading"
                    @click="delTask(slotProps.data.id)"
                />
            </template>
        </Column>
        <template #footer>
            In total there are {{ products ? products.length : 0 }} tasks.
        </template>
    </DataTable>
    <loadingScreen v-if="tasks.isLoading" />
</template>
