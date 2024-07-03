https://login.microsoftonline.com/{tenant_id}/oauth2/v2.0/authorize?client_id={client_id}&response_type=code&redirect_uri=https://localhost&response_mode=query&scope=https://graph.microsoft.com/.default&state=12345


Content-Type: application/x-www-form-urlencoded


grant_type: authorization_code
client_id: {client_id} (your application/client ID from Azure AD)
client_secret: {client_secret} (the client secret you generated)
code: {authorization_code} (the authorization code you just obtained)
redirect_uri: https://localhost (or your specified redirect URI)
scope: https://graph.microsoft.com/.default
