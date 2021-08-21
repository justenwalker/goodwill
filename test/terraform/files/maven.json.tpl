{
  "repositories": [
    {
      "id": "sonatype-staging",
      "layout": "default",
      "auth": {
        "username": "${sonatype_username}",
        "password": "${sonatype_password}"
      },
      "url": "https://s01.oss.sonatype.org/content/repositories/${sonatype_staging_repo}/"
    },
    {
      "id": "sonatype-public",
      "layout": "default",
      "url": "https://s01.oss.sonatype.org/content/groups/public/"
    }
  ]
}
