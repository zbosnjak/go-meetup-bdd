Go-meetup-bdd
========
Small app for managing bank accounts. 

# App Runner

### Build the app:

       for windows:
          $  make build GOOS=windows    
       for linux:
          $  make build GOOS=linux    
       for mac (default):
          $  make build GOOS=darwin    
          
### Run the app:

       using default port 9099:
          $  ./bank_acc    
       using specific port:
          $  ./bank_acc rest_port=9090    

NOTE: If app is not started at '127.0.0.1:9099', to successfully run tests change env variable in 'assets/envs/test'. 
      
# Test Runner

### Run the full suite:

       for one feature:
          $  make tests-run 

### Test one feature:

       for one feature:
          $  make tests-run FEATURE=S_100_bank_acc

### Run tests and generate result files:

     for the full suite:
        $  make tests-run FULL_REPORT=1

     for one feature:
        $  make tests-run REPORT=S_100_bank_acc


### Compile results into an html report:
NOTE: Before running this command remove any log messages from all .json files in 'tests/results'

        $ make tests-report