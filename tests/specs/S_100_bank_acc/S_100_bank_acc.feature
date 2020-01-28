Feature: Transferring money between accounts
    In order to manage my money more efficiently
    As a bank client
    I want to transfer funds between my accounts whenever I need to

Background:
    Given I want to manage my accounts
    And transfer money from one to another

Scenario: Transferring money to a savings account
    Given my Current account has a balance of 1000.00
    And my Savings account has a balance of 2000.00
    When I transfer 500.00 from my Current account to my Savings account
    Then I should have 500.00 in my Current account
    And I should have 2500.00 in my Savings account

Scenario: Transferring with insufficient funds
    Given my Current account has a balance of 1000.00
    And my Savings account has a balance of 2000.00
    When I transfer 1500.00 from my Current account to my Savings account
    Then I should receive an "Insufficient funds" error
    Then I should have 1000.00 in my Current account
    And I should have 2000.00 in my Savings account

Scenario Outline: Earning interest
    Given I have a saving account with a balance of <initial-balance>
    When the monthly interest <interest-rate> is calculated
    Then I should have <new-balance> at end of year
    Examples:
    | initial-balance  | interest-rate | new-balance |
    | 10000 |   1 | 11268,249 |
    | 10000 |   3 | 14257,6045 |
    | 10000 |  5 | 17958,55 |

Scenario: Multiple transaction
    Given my Current account has a balance of 1000.00
    And my Savings account has a balance of 2000.00
    When I perform <action> with <amount>
      | action          | amount |
      | transfer        | 100    |
      | deposit         | 1000   |
      | deposit         | 500    |
      | transfer        | 1000   |
      | transfer        | 10000  |
    Then I should have 1400.00 in my Current account
        And I should have 3100.00 in my Savings account

@manual
Scenario: Scenario to hard to automate
    Given I have a scenario hard to automate
    And I am too lazy
    When I put manual tag
    Then scenario will be skipped when running tests from this feature file