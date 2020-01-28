var reporter = require('cucumber-html-reporter');

var options = {
    brandTitle: "GOlang Meetup",
    name: "Bank Account App Integration Test Suite",
    theme: 'bootstrap',
    jsonDir: '/test/results',
    output: '/test/reports/report.html',
    reportSuiteAsScenarios: true,
    launchReport: false,
    metadata: {
        "Test Environment": "Docker",
    }
};

reporter.generate(options);