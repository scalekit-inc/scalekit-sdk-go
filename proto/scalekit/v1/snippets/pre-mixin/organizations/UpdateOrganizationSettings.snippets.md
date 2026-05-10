---
operationId: OrganizationService_UpdateOrganizationSettings
---

```javascript
const settings = {
  features: [
    {
      name: 'sso',
      enabled: true,
    },
    {
      name: 'dir_sync',
      enabled: true,
    },
  ],
};

await scalekit.organization.updateOrganizationSettings('<organization_id>', settings);
```

```python
settings = [
        {
            "name": "sso",
            "enabled": True
        },
        {
            "name": "dir_sync",
            "enabled": True
        }
    ]

scalekit_client.organization.update_organization_settings(
  organization_id='<organization_id>', settings=settings
)
```

## Go SDK

```go
settings := OrganizationSettings{
		Features: []Feature{
			{
				Name:    "sso",
				Enabled: true,
			},
			{
				Name:    "dir_sync",
				Enabled: true,
			},
		},
	}

organization,err := scalekitClient.Organization().UpdateOrganizationSettings(ctx, organizationId, settings)
```

## Java SDK

```java
OrganizationSettingsFeature featureSSO = OrganizationSettingsFeature.newBuilder()
                .setName("sso")
                .setEnabled(true)
                .build();

OrganizationSettingsFeature featureDirectorySync = OrganizationSettingsFeature.newBuilder()
                .setName("dir_sync")
                .setEnabled(true)
                .build();

updatedOrganization = scalekitClient.organizations()
                .updateOrganizationSettings(organization.getId(), List.of(featureSSO, featureDirectorySync));
```
