<template>
    <div class="c-horizontal">
        <div class="c-left">
            <slot name="left"></slot>
        </div>

        <div :class="['c-main', { 'h-center': centerContent }]">
            <slot></slot>
        </div>

        <div class="c-right">
            <slot name="right"></slot>
        </div>
    </div>
</template>

<script setup>
defineProps({
    centerContent: {
        type: Boolean,
        default: false
    }
})
</script>

<style>
.c-horizontal {
    display: grid;
    grid-template-columns: auto 1fr auto;
    width: 100%;
}
.c-horizontal > .c-main.h-center {
    display: grid;
    place-items: normal center;
}
</style>
