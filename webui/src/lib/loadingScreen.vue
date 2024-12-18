<script setup>
import ProgressSpinner from 'primevue/progressspinner'
</script>
<style>
.c-loading {
    display: flex;
    background-color: rgba(133, 133, 133, 0.5);

    position: fixed;
    height: 100%;
    width: 100%;
    left: 0;
    top: 0;

    justify-content: center;
    align-items: center;
    pointer-events: auto;
    z-index: 1101;
    opacity: 1;
    transition: opacity 1s linear;
}
</style>

<template>
    <div class="c-loading p-dialog-mask p-component-overlay" data-pc-section="mask">
        <span
            class="p-hidden-accessible p-hidden-focusable"
            tabindex="0"
            role="presentation"
            aria-hidden="true"
            data-p-hidden-accessible="true"
            data-p-hidden-focusable="true"
            data-pc-section="firstfocusableelement"
        ></span>

        <ProgressSpinner aria-label="Loading" />
    </div>
</template>
