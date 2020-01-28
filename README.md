Go-meetup-bdd
========
Small app for managing bank accounts. 



# Test Runner

## Examples

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
        $  make tests-run REPORT=S-03581_login


### Compile results into an html report:
NOTE: Before running this command remove any log messages from all .json files in 'tests/results'

        $ make tests-report