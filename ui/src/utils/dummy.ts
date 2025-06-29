import type {ModelsDailyTraffic, ModelsIPBan, ModelsOnlineUserSession} from "@/api";

const dummyTrafficData = <ModelsDailyTraffic[]>([
    {date: '2025-06-18', rx: 1.2, tx: 2.5},
    {date: '2025-06-19', rx: 0.9, tx: 1.1},
    {date: '2025-06-21', rx: 0.7, tx: 0.8},
    {date: '2025-06-22', rx: 1.0, tx: 1.3},
    {date: '2025-06-23', rx: 0.5, tx: 0.6},
    {date: '2025-06-25', rx: 0.3, tx: 0.4},
    {date: '2025-06-26', rx: 1.5, tx: 2.0},
    {date: '2025-06-27', rx: 2.1, tx: 3.2},
    {date: '2025-06-28', rx: 10, tx: 4.0},
])

const dummyOnlineUsers = <Array<ModelsOnlineUserSession>>([
    {
        "Username": "masoud1",
        "Groupname": "(none)",
        "Average RX": "12.3 kB/s",
        "Average TX": "1.2 kB/s",
        "_Connected at": "20s"
    },
    {
        "Username": "jane_doe",
        "Groupname": "group_test",
        "Average RX": "34.6 kB/s",
        "Average TX": "5.7 kB/s",
        "_Connected at": "65m:20s"
    },
    {
        "Username": "admin",
        "Groupname": "group_test2",
        "Average RX": "98.1 kB/s",
        "Average TX": "22.4 kB/s",
        "_Connected at": "1h:30m:40s"
    }
])

const dummyBanIPs = <Array<ModelsIPBan>>([
    {
        "IP": "172.17.0.1",
        "Since": "2025-06-28 18:26",
        "_Since": " 4m:55s",
        "Score": 80
    },
    {
        "IP": "172.17.0.2",
        "Since": "2025-06-28 18:26",
        "_Since": " 9m:55s",
        "Score": 120
    },
    {
        "IP": "172.17.0.3",
        "Since": "2025-06-28 19:26",
        "_Since": " 10m:55s",
        "Score": 160
    },
    {
        "IP": "172.17.0.4",
        "Since": "2025-06-29 23:26",
        "_Since": " 1h:10m:55s",
        "Score": 220
    },
    {
        "IP": "172.17.0.5",
        "Since": "2025-06-31 23:26",
        "_Since": " 1h:10m:55s",
        "Score": 32
    },
    {
        "IP": "172.17.0.6",
        "Since": "2025-06-29 23:26",
        "_Since": " 1h:10m:55s",
        "Score": 190
    }
])

export {
    dummyTrafficData,
    dummyOnlineUsers,
    dummyBanIPs
}