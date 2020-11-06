# Cadence worker sample
This is a simple example on how to write and run Cadence workflows and activities

Dependencies:
- go 1.13+
- docker
- docker-compose

1. Clone cadence:
```
export CADENCE_REPO_LOCATION='< your-desired-cadence-repo-location >'
git clone git@github.com:uber/cadence.git $CADENCE_REPO_LOCATION`
```

2. Start cadence:
```
cd $CADENCE_REPO_LOCATION/docker
docker-compose up
```

3. Build CLI:
```
cd $CADENCE_REPO_LOCATION
make bins
```

4. Register domain:
```
$CADENCE_REPO_LOCATION/cadence --domain simple-domain domain register --global_domain false
```

5. Run worker:
```
cd '< this-repo-location >'
go build ./...
./cadence-worker
```

6. Execute branch workflow:
```
$CADENCE_REPO_LOCATION/cadence --domain simple-domain workflow run --tasklist simple-worker --workflow_type main.sampleBranchWorkflow --execution_timeout 5
```

7. Execute parallel workflow:
```
$CADENCE_REPO_LOCATION/cadence --domain simple-domain workflow run --tasklist simple-worker --workflow_type main.sampleParallelWorkflow --execution_timeout 5
```

8. Checkout the workflow status in Cadence UI: http://localhost:8088/domains/simple-domain/workflows?range=last-3-days&status=COMPLETED