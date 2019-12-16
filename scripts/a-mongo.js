db.createUser(
    {
        user: "paul",
        pwd: "paul123",
        roles: [
            {
                role: "readWrite",
                db: "cdr"
            }
        ]
    }
)