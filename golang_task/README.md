# Golang ToDO List CLI Tool 🛠️


## Description
Mytask is a straightforward command-line tool designed to help you manage your daily tasks efficiently. Built with Go, it offers a user-friendly interface to create, view, update, and delete tasks, all stored in a local JSON file.

![Project Image](https://github.com/dev-dhanushkumar/Golang-Projects/blob/main/golang_task/mytask_list.png)


## Table of Content
- [Description](#description)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#liclicense)


## Installation
Mytask is a breeze to install! Here's how to get started:

### Prerequisites:

Go: Ensure you have Go installed on your system. You can download it from the official website: https://golang.org/dl/

### Steps:
1. **Clone the Repository**:
Open a terminal or command prompt and navigate to your desired project directory. Then, clone the GoTodo repository using Git:
    ```bash
    git clone https://github.com/dev-dhanushkumar/Golang-Projects.git
    ```
2. **Change Directory**:
Navigate to the project root directory:
    ```bash
    cd Golang-Projects/golang_task  
    ```
3. **Build the Project**:
Compile the Go source code to create the executable file:
    ```bash
    go build
    ```
    - **Windows**: This will typically generate a file named `mytask.exe`.
    - **Linux/macOS**: It will usually create a file named `mytask`.
4. **Configure Environment Variables (Optional, but Recommended):**
    - To create this folder in your C derive `c:\mytask\bin\` and in this bin folder to paste that project mytask.exe file.
    - Adding the executable location to your system's environment path allows you to run mytask from any directory:

        **Windows:**
        - Search for "Environment Variables" in your system settings.
        - Click on "Edit the system environment variables".
        - Under "User variables" or "System variables" (depending on your preference), find the "Path" variable and click "Edit".
        - Click "New" and add the directory containing the mytask.exe file (C:\mytask\bin).
        - Click "OK" on all open windows to save the changes.

        **Linux/macOS:**
        - Open terminal in from build path after to change that build to executable so execute below command
            ```bash
            sudo chmod +x mytask
            ```
        - After move `mytask` file to `/usr/local/bin` so to execute this below command,
            ```bash
            sudo mv mytask /usr/local/bin/
            ```
5. **Verify Installation:**
    - **On Windows:** Open Command Prompt and type the following to check if the installation was successful:
        ```bash
        mytask help
        ```
    - **On Linux:** Open a new terminal window and type the following to check if the installation was successful:
        ```bash
        mytask help
        ```
    If everything is set up correctly, this command will display the help message for the mytask application.


## Usage

### Initializing GoTodo:
Before you start using mytask, you need to initialize it to create a JSON file to store your tasks. Run the following command
```bash
mytask init #This will create a .gtodo.json file in your home directory.
```
Be sure to `mytask init` to generate an empty JSON file in your home directory to store todo tasks.

### Adding a Task:
To add a new task, use the add command followed by the task description and optional category:
```bash
mytask add -task "Implement login feature with JWT authentication " -cat "Feature"
```
### Listing Tasks:
To list all tasks, use the `list` command:
```bash
mytask list
```
To filter tasks based on completion status or category, use the following options:
```bash
mytask list -done 1  # List completed tasks
mytask list -cat "Work"  # List tasks in the "Work" category
mytask list -done 1 -cat "Work"  # List completed tasks in the "Work" category
```

### Updating a Task:
To update an existing task, use the update command followed by the task ID, new task description, and optional new category:
```bash
mytask update -id 1 -task "Finish the report" -cat "Work" # Here 1 - taskID_number
```
Here `-done 1` means task completed.
To change done status of task use below command.
```bash
mytask update -id 1 -done 1
```
### Deleting a Task:
To delete a task, use the `delete` command followed by the task ID:
```bash
mytask delete -id 1
```

## Contributing

If you find a bug or have a feature request, please open an issue on the GitHub repository. Pull requests are also welcome!

## License

This application is licensed under the MIT License. See the LICENSE file for details.
