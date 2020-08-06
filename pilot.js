const admin = require("firebase-admin")
const readline = require("readline")

admin.initializeApp({
  credential: admin.credential.applicationDefault(),
  databaseURL: "https://chariot-rc.firebaseio.com",
})

const db = admin.database()
const ref = db.ref("chariot_1")

readline.emitKeypressEvents(process.stdin)
process.stdin.setRawMode(true)
process.stdin.on("keypress", async (str, key) => {
  let engine = 0
  if (key.name === "w") {
    engine = 1
  } else if (key.name === "s") {
    engine = -1
  } else if (key.name === "escape") {
    process.exit(0)
  } else {
    engine = 0
  }

  await ref.set({ engine: engine })
  console.log(`engine: ${engine}`)
})

console.log("Press W or S")
