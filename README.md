## Set up postgres

This project need postgres database to use. Currently in development state the set up involve starting postgresql as a services via homebrew.

```bash
brew services start postgres
```

```bash
psql -d greenlight
```

If this not work, it probably because the greenlight database haven't created yet, we would need to connect as default user and create this `greenlight` database manually.
