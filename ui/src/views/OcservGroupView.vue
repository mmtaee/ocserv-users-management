<script lang="ts" setup>
// import type {DataTableHeader} from "vuetify/framework";
import {useLocale} from "vuetify/framework";
import {reactive, ref} from 'vue'
import {numberRule} from "@/utils/rules.ts";

const {t} = useLocale()

// const items = [
//   {uid: "01ARZ3NDEKTSV4RRFFQ69G5FAV", name: "Anc 1234"},
//   {uid: "01ARZ3NDEKTSV4RRFFQ69G5FAV", name: "Anc 4568"},
//   {uid: "01ARZ3NDEKTSV4RRFFQ69G5FAV", name: "Anc 1248"},
//   {uid: "01ARZ3NDEKTSV4RRFFQ69G5FAV", name: "Anc 1298"},
// ]
//
// const headers: DataTableHeader[] = [
//   {title: t("UID"), key: 'uid', align: 'start', sortable: false},
//   {title: t("Name"), key: 'name', align: 'start'},
//   {title: t("ACTIONS"), key: 'actions', align: 'center', sortable: false},
// ]


const tab = ref('option-1')

// const valid = ref(true)


const fields = [
  // Text Fields
  {key: 'dns', label: 'DNS (comma-separated)', type: 'text', hint: 'e.g., 8.8.8.8,1.1.1.1'},
  {key: 'nbns', label: 'NBNS', type: 'text', hint: 'NetBIOS Name Server (e.g., 192.168.1.10)'},
  {key: 'ipv4Network', label: 'IPv4 Network', type: 'text', hint: 'CIDR notation (e.g., 192.168.0.0/24)'},
  {key: 'explicitIPv4', label: 'Explicit IPv4', type: 'text', hint: 'Specific IP address (e.g., 192.168.1.5)'},
  {key: 'cgroup', label: 'CGroup', type: 'text', hint: 'Linux control group name (e.g., net_cls)'},
  {key: 'iroute', label: 'Internal Route', type: 'text', hint: 'Custom internal route (e.g., 10.0.0.0/8)'},
  {
    key: 'route',
    label: 'Route (comma-separated)',
    type: 'text',
    hint: 'Routes to push (e.g., 10.0.0.0/8,192.168.1.0/24)'
  },
  {key: 'noRoute', label: 'No Route (comma-separated)', type: 'text', hint: 'Routes to exclude (e.g., 172.16.0.0/12)'},
  {key: 'restrictUserToPorts', label: 'Restrict User To Ports', type: 'text', hint: 'Allowed ports (e.g., 80,443)'},
  {
    key: 'splitDNS',
    label: 'Split DNS (comma-separated)',
    type: 'text',
    hint: 'Domain list for split DNS (e.g., internal.local,corp.net)'
  },

  // Number Fields
  {
    key: 'rxDataPerSec',
    label: 'RX Data Per Sec',
    type: 'number',
    hint: 'Max receive speed in bytes/sec (e.g., 100000)'
  },
  {
    key: 'txDataPerSec',
    label: 'TX Data Per Sec',
    type: 'number',
    hint: 'Max transmit speed in bytes/sec (e.g., 100000)'
  },
  {key: 'netPriority', label: 'Net Priority', type: 'number', hint: 'Traffic class priority (0 = default)'},
  {key: 'keepAlive', label: 'KeepAlive', type: 'number', hint: 'Interval to send keepalive packets (seconds)'},
  {key: 'dpd', label: 'DPD Timeout', type: 'number', hint: 'Dead Peer Detection timeout (seconds)'},
  {key: 'mobileDPD', label: 'Mobile DPD Timeout', type: 'number', hint: 'DPD timeout for mobile clients (seconds)'},
  {
    key: 'maxSameClients',
    label: 'Max Same Clients',
    type: 'number',
    hint: 'Max concurrent logins with the same username'
  },
  {key: 'statsReportTime', label: 'Stats Report Time', type: 'number', hint: 'Interval for reporting stats (seconds)'},
  {key: 'mtu', label: 'MTU', type: 'number', hint: 'Max transmission unit (e.g., 1500)'},
  {key: 'idleTimeout', label: 'Idle Timeout', type: 'number', hint: 'Disconnect after inactivity (seconds)'},
  {
    key: 'mobileIdleTimeout',
    label: 'Mobile Idle Timeout',
    type: 'number',
    hint: 'Inactivity timeout for mobile users (seconds)'
  },
  {key: 'sessionTimeout', label: 'Session Timeout', type: 'number', hint: 'Max session duration (seconds)'},

  // Boolean Switches
  {key: 'denyRoaming', label: 'Deny Roaming', type: 'switch', hint: 'Prevent client session handover between IPs'},
  {key: 'noUDP', label: 'Disable UDP', type: 'switch', hint: 'Force TCP-only connections (no UDP transport)'},
  {key: 'tunnelAllDNS', label: 'Tunnel All DNS', type: 'switch', hint: 'Send all DNS queries through VPN tunnel'},
  {
    key: 'restrictUserToRoutes',
    label: 'Restrict User To Routes',
    type: 'switch',
    hint: 'Limit user to configured route list'
  }
]


interface FormFields {
  dns: string;
  nbns: string;
  ipv4Network: string;
  explicitIPv4: string;
  cgroup: string;
  iroute: string;
  route: string;
  noRoute: string;
  restrictUserToPorts: string;
  splitDNS: string;

  rxDataPerSec: number;
  txDataPerSec: number;
  netPriority: number;
  keepAlive: number;
  dpd: number;
  mobileDPD: number;
  maxSameClients: number;
  statsReportTime: number;
  mtu: number;
  idleTimeout: number;
  mobileIdleTimeout: number;
  sessionTimeout: number;

  denyRoaming: boolean;
  noUDP: boolean;
  tunnelAllDNS: boolean;
  restrictUserToRoutes: boolean;
}

const form = reactive<FormFields>({
  dns: '8.8.8.8,1.1.1.1',
  nbns: '192.168.1.1',
  ipv4Network: '192.168.1.0/24',
  rxDataPerSec: 100000,
  txDataPerSec: 200000,
  explicitIPv4: '192.168.100.10',
  cgroup: 'cpuset,cpu:test',
  iroute: '10.0.0.0/8',
  route: '0.0.0.0/0,10.10.0.0/16',
  noRoute: '192.168.0.0/16,10.0.0.0/8',
  netPriority: 1,
  denyRoaming: true,
  noUDP: false,
  keepAlive: 60,
  dpd: 90,
  mobileDPD: 300,
  maxSameClients: 2,
  tunnelAllDNS: true,
  statsReportTime: 300,
  mtu: 1400,
  idleTimeout: 600,
  mobileIdleTimeout: 900,
  restrictUserToRoutes: true,
  restrictUserToPorts: 'tcp(443),tcp(80),udp(53)',
  splitDNS: 'example.com,internal.company.com',
  sessionTimeout: 3600,
})

</script>

<template>
  <v-row>
    <v-col>
      <v-card min-height="850">
        <v-toolbar :title="t('OCSERV_GROUPS')" color="primary">
        </v-toolbar>

        <div class="d-flex flex-row">
          <v-tabs
              v-model="tab"
              color="primary"
              direction="vertical"
          >
            <v-tab :text="t('DEFAULTS')" value="defaults"></v-tab>
            <v-tab :text="t('OTHER')" value="other"></v-tab>
          </v-tabs>

          <v-tabs-window v-model="tab">
            <v-tabs-window-item value="defaults">
              <v-card flat>
                <v-card-title class="text-end">
                  <v-btn color="primary" variant="outlined">
                    {{ t("SAVE") }}
                  </v-btn>
                </v-card-title>
                <v-card-text>
                  <v-form>
                    <v-row dense>

                      <v-col cols="12">
                        <h3>Text Fields</h3>
                        <v-divider/>
                      </v-col>

                      <template v-for="field in fields.filter(f => f.type === 'text')" :key="field.key">
                        <v-col
                            class="my-1"
                            cols="12"
                            lg="3"
                            md="4"
                            xl="2"
                        >
                          <v-text-field
                              v-model="form[field.key as keyof typeof form]"
                              :hint="field.hint"
                              :label="field.label"
                              density="default"
                              persistent-hint
                              type="text"
                              variant="underlined"
                          />
                        </v-col>
                      </template>

                      <!-- Number Fields Section -->
                      <v-col class="mt-6" cols="12">
                        <h3>Number Fields</h3>
                        <v-divider/>
                      </v-col>

                      <template v-for="field in fields.filter(f => f.type === 'number')" :key="field.key">
                        <v-col
                            class="my-1"
                            cols="12"
                            lg="3"
                            md="4"
                            xl="2"
                        >
                          <v-text-field
                              v-model.number="form[field.key as keyof typeof form]"
                              :hint="field.hint"
                              :label="field.label"
                              :min="0"
                              controlVariant="hidden"
                              density="default"
                              persistent-hint
                              variant="underlined"
                          />
                        </v-col>
                      </template>

                      <!-- Switch Fields Section -->
                      <v-col class="mt-6" cols="12">
                        <h3>Options</h3>
                        <v-divider/>
                      </v-col>

                      <template v-for="field in fields.filter(f => f.type === 'switch')" :key="field.key">
                        <v-col
                            class="my-1"
                            cols="12"
                            md="2"
                        >
                          <v-switch
                              v-model="form[field.key as keyof typeof form]"
                              :label="field.label"
                              class="ms-1"
                              color="primary"
                          />
                        </v-col>
                      </template>
                    </v-row>
                  </v-form>

                </v-card-text>

              </v-card>
            </v-tabs-window-item>

            <v-tabs-window-item value="other">
              <v-card flat>
                <v-card-text>
                  <p>
                    Morbi nec metus. Suspendisse faucibus, nunc et pellentesque egestas, lacus ante convallis tellus,
                    vitae iaculis lacus elit id tortor. Sed mollis, eros et ultrices tempus, mauris ipsum aliquam
                    libero,
                    non adipiscing dolor urna a orci. Curabitur ligula sapien, tincidunt non, euismod vitae, posuere
                    imperdiet, leo. Nunc sed turpis.
                  </p>

                  <p>
                    Suspendisse feugiat. Suspendisse faucibus, nunc et pellentesque egestas, lacus ante convallis
                    tellus,
                    vitae iaculis lacus elit id tortor. Proin viverra, ligula sit amet ultrices semper, ligula arcu
                    tristique sapien, a accumsan nisi mauris ac eros. In hac habitasse platea dictumst. Fusce ac felis
                    sit
                    amet ligula pharetra condimentum.
                  </p>

                  <p>
                    Sed consequat, leo eget bibendum sodales, augue velit cursus nunc, quis gravida magna mi a libero.
                    Nam
                    commodo suscipit quam. In consectetuer turpis ut velit. Sed cursus turpis vitae tortor. Aliquam eu
                    nunc.
                  </p>

                  <p>
                    Etiam ut purus mattis mauris sodales aliquam. Ut varius tincidunt libero. Aenean viverra rhoncus
                    pede.
                    Duis leo. Fusce fermentum odio nec arcu.
                  </p>

                  <p class="mb-0">
                    Donec venenatis vulputate lorem. Aenean viverra rhoncus pede. In dui magna, posuere eget, vestibulum
                    et, tempor auctor, justo. Fusce commodo aliquam arcu. Suspendisse enim turpis, dictum sed, iaculis
                    a,
                    condimentum nec, nisi.
                  </p>
                </v-card-text>
              </v-card>
            </v-tabs-window-item>

          </v-tabs-window>
        </div>
      </v-card>
    </v-col>
  </v-row>


  <!--  <v-row align="start" justify="center">-->
  <!--    <v-col cols="12" lg="7" md="9" sm="12">-->
  <!--      <v-card>-->
  <!--        <v-card-text>-->
  <!--          <v-data-table-->
  <!--              :headers="headers"-->
  <!--              :items="items"-->
  <!--              hide-default-footer-->
  <!--          >-->
  <!--            <template v-slot:top>-->
  <!--              <v-toolbar flat>-->
  <!--                <v-toolbar-title class="text-capitalize">-->
  <!--                  <v-icon class="mb-3" color="medium-emphasis" icon="mdi-router-network" size="large" start></v-icon>-->
  <!--                  {{ t("OCSERV_GROUPS") }}-->
  <!--                </v-toolbar-title>-->

  <!--                <v-btn-->
  <!--                    :text="t('CREATE')"-->
  <!--                    border-->
  <!--                    class="me-2"-->
  <!--                    prepend-icon="mdi-plus"-->
  <!--                    rounded="lg"-->
  <!--                />-->
  <!--              </v-toolbar>-->
  <!--            </template>-->

  <!--            <template #item.actions="{ }">-->
  <!--              <v-menu>-->
  <!--                <template v-slot:activator="{ props }">-->
  <!--                  <v-btn icon="mdi-dots-vertical" v-bind="props" variant="text"></v-btn>-->
  <!--                </template>-->

  <!--                <v-list>-->
  <!--                  <v-list-item>-->
  <!--                    <v-icon class="mb-1" color="info" start>mdi-pencil</v-icon>-->
  <!--                    <span class="text-info">{{ t('CONFIG') }}</span>-->
  <!--                  </v-list-item>-->


  <!--                  <v-list-item>-->
  <!--                    <v-icon class="mb-1" color="info" start>mdi-pencil</v-icon>-->
  <!--                    <span class="text-info">{{ t('EDIT') }}</span>-->
  <!--                  </v-list-item>-->

  <!--                  <v-list-item>-->
  <!--                    <v-icon class="mb-1" color="red" start>mdi-delete</v-icon>-->
  <!--                    <span class="text-error">{{ t('DELETE') }}</span>-->
  <!--                  </v-list-item>-->
  <!--                </v-list>-->
  <!--              </v-menu>-->
  <!--            </template>-->

  <!--          </v-data-table>-->
  <!--        </v-card-text>-->
  <!--      </v-card>-->
  <!--    </v-col>-->

  <!--  </v-row>-->

</template>
