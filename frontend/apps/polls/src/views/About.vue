<template>
    <div class="about">
        <h1>This is an about page</h1>
        <form @submit.prevent="greet">
            <InputText type="text" v-model="text" />
            <Button type="submit" label="Submit" />
            <h3>{{ pmessage }}</h3>
        </form>
        <DataView :value="polls" layout="grid">
            <template #grid="slotProps">
                <div class="p-col-12 p-md-4">
                    <div class="product-grid-item card">
                        <div class="product-grid-item-top">
                            <div>
                                <i class="pi pi-tag product-category-icon"></i>
                                <span class="product-category">{{
                                    slotProps.data.title
                                }}</span>
                            </div>
                        </div>
                        <div class="product-grid-item-bottom">
                            <Button icon="pi pi-shopping-cart"></Button>
                        </div>
                    </div>
                </div>
            </template>
        </DataView>
    </div>
</template>

<script lang="ts">
import { Vue } from 'vue-class-component'
import store from '../store'

export default class About extends Vue {
    public message = ''
    public text = ''
    public polls = []

    mounted() {
        store.dispatch('bindPolls').then(() => (this.polls = store.state.polls))
    }

    unmounted() {
        store.dispatch('unbindPolls').then(() => {
            console.log('terminated db connection')
        })
    }

    public greet() {
        if (this.text.length) {
            this.message = this.text
        }
    }
}
</script>
