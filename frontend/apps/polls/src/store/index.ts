import { firestoreAction, vuexfireMutations } from 'vuexfire'
import { createStore } from 'vuex'
import { db } from '../db'

export default createStore({
    state: {
        polls: [],
        test: ['test'],
    },
    mutations: {
        ...vuexfireMutations,
    },
    actions: {
        bindPolls: firestoreAction(({ bindFirestoreRef }) => {
            return bindFirestoreRef('polls', db.collection('polls'))
        }),
        unbindPolls: firestoreAction(({ unbindFirestoreRef }) => {
            unbindFirestoreRef('polls')
        }),
    },
    modules: {},
})
