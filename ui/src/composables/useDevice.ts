import {computed, ref} from 'vue';
import {useDisplay} from 'vuetify';

export const isSmallDevice = ref<boolean>(false);

export function useIsSmallDevice(): boolean {
    const {xs, sm, mobile} = useDisplay();
    const isSmall = computed(() => mobile.value || xs.value || sm.value).value;
    isSmallDevice.value = isSmall;
    return isSmall;
}
