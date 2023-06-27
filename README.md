# Simple Budget Application
*Mixed Machine* <br />
*mixedmachine.dev@gmail.com*

## Description
This is a simple budget application that allows the user to add expenses and deposits to their budget then allocate the funds to track where the income is going. The application is downloadable and uses sqlite for the database. In the future, I would like to add a login feature and allow the user to save their budget to a remote database to access the same information on different devices. This will include adding users to the current models. I would also like to add a feature to allow the user to export and import snapshots of different months. This would allow the user to have a budget for each month and be able to look at previous months. I also need to figure out how to highlight expenses based on what is allocated.


## Table of Contents
* [Installation](#installation)
* [Usage](#usage)
* [Future Features](#future-features)
* [Bugs](#bugs)
* [License](#license)


## Installation
To install the application, clone the store and run `make init` or `make build.win/lin` or `make run` to install the dependencies.

If you do not have make installed, you can run `go mod download` to install the dependencies. The application can be run with `go run main.go` or `go build main.go` and then `./main.exe` or `./main` depending on your operating system.

## Usage
![image](./pictures/preview1.png)
![image](./pictures/preview2.png)
![image](./pictures/preview3.png)
![image](./pictures/preview4.png)
![image](./pictures/preview5.png)
![image](./pictures/preview6.png)


## Future Features
- [ ] Expense highlighting based on allocated funds
    - [ ] Red if nothing allocated
    - [ ] Yellow if partially allocated
    - [ ] Green if fully allocated
- [x] Allocation of each Income
- [x] Auto-allocation based on available income and date of expenses
- [ ] Button to convert all dates of budget to next month
- [ ] Button to convert all dates of budget to previous month
- [ ] Button to convert all dates of budget to current month
    - [ ] With these features I do not want to save each month's budget, but rather have the user be able to change the month they are looking at for convenience
    - [ ] Due to the above, there should be an export button to "snapshot" the current budget
- [ ] Add a "snapshot" button to save the current budget
- [ ] Add a "load" button to load a saved budget
- [x] Add a "clear" button to clear the current budget
- [x] enable the user to add notes
- [ ] enable a switch to save locally or to remote database (paid feature or bring your own database)
- [ ] Add login (paid feature)
- [ ] Search on expenses & allocations
- [ ] Add sort by button on expenses & allocations

## Bugs
- [ ] The application uses local storage to save the budget items
- [ ] Allocations of income do not delete when the income is deleted
- [ ] On mobile, backspace double deletes
- [x] In light mode, the headers are not visible
- [x] Crashes when the user tries to update an expense name
- [x] Text hint on allocations edit is wrong
- [x] Allocation doesn't need to do anything to income.Allocation
- [x] income.Allocation should be removed since switching to sql database for getting income allocaitons
- [ ] Icon doesn't show up with installer
- [ ] Using the same name for income and expense causes the application to crash

\* Add more bugs in github issues as they're found


## License:
This project is licensed under the MIT License - see the 
[LICENSE.md](./LICENSE.txt) file for details.
