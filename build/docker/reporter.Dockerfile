FROM node:latest

RUN npm install cucumber-html-reporter@4.0.3 --save-dev

CMD node /test/reports/generate-report.js