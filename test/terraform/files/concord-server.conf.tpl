concord-server {
    port = 8001
    db {
        url = "jdbc:postgresql://db:5432/postgres"
        appPassword = "${pg_password}"
        inventoryPassword = "${pg_password}"
    }
    secretStore {
        serverPassword = "cTFxMXExcTE="
        secretStoreSalt = "SCk4KmBlazMi"
        projectSecretSalt = "I34xCmcOCwVv"
    }
    ldap {
        url = "ldap://ldap:389"
        searchBase = "ou=Users,${ldap_base_dn}"
        principalSearchFilter = "(cn={0})"
        userSearchFilter = "(cn={0})"
        usernameProperty = "cn"
        userPrincipalNameProperty = ""
        returningAttributes = ["*", "memberOf"]
        systemUsername = "cn=admin,${ldap_base_dn}"
        systemPassword = "${ldap_admin_password}"
    }
}