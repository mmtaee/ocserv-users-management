import {defineStore} from 'pinia';

export const useIsSmallDisplay = defineStore('display', {
    state: () => ({
        isSmall: false,
    }),
    actions: {
        setIsSmall(val: boolean) {
            this.isSmall = val;
        },
    },
    getters: {
        isSmallDisplay: (state) => {
            return state.isSmall
        }
    }
});