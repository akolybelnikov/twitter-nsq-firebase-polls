import firebase from 'firebase/app'
import 'firebase/firestore'

export const db = firebase
    .initializeApp({ projectId: process.env.VUE_APP_FIREBASE_PROJECT_ID })
    .firestore()
