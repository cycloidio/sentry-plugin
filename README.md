# Sentry Plugin

This Plugin fetches all the information from Sentry into the Plugin DB

## Configuration

### sentry_api_key

The Key to use on the API

### sentry_endpoint

If it's a different endpoint than `https://sentry.io/api/0`

### sentry_organization_slug

If you don't want to import all the Organizations then set this to the Slug of the org

## Test

To test this we need access to Sentry, for it you have to have a `.env` with `SENTRY_API_KEY="..."` in it
