const admin = require("firebase-admin")

admin.initializeApp({
  credential: admin.credential.applicationDefault(),
  databaseURL: "https://chariot-rc.firebaseio.com",
})

const db = admin.database()
const ref = db.ref("chariot_1")

ref.on("child_changed", (snapshot) => {
  console.log(snapshot.val())
})
