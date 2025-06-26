import {computed, ref} from 'vue';
import {useDisplay} from 'vuetify';

export const isSmallDevice = ref<boolean>(false);

export function useIsSmallDevice(): boolean {
    const display = useDisplay()
    const isSmallDisplay = computed(() => display.mdAndDown.value)

    isSmallDevice.value = isSmallDisplay.value;
    return isSmallDisplay.value;
}
