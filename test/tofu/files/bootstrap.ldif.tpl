# == Users
dn: ou=Users,${ldap_base_dn}
objectClass: top
objectClass: organizationalUnit
ou: Users

## User: concord-admin
dn: cn=concord-admin,ou=Users,${ldap_base_dn}
cn: concord-admin
objectClass: top
objectClass: simpleSecurityObject
objectClass: organizationalRole
description: Concord Administrator
userPassword: {SSHA}IY+YzW3SyzL1LRFMaHm5x1SFcDXxTmoC

# == Groups
dn: ou=Groups,${ldap_base_dn}
objectClass: top
objectClass: organizationalUnit
ou: Groups

# Group: concord-admins
dn: cn=concord-admins,ou=Groups,${ldap_base_dn}
cn: concord-admins
description: Concord Admins
objectClass: top
objectClass: groupOfUniqueNames
uniqueMember: cn=concord-admin,ou=Users,${ldap_base_dn}