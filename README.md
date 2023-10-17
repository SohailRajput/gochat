# Socket-Based Chat App in GoLang

This is a simple command line based chat application built in GoLang that allows users to communicate in real-time over a local network. The app is designed for learning purposes and showcases the basics of socket programming in Go.

## Getting Started

To run the chat application, follow these steps:

1. Build the `app.go` file.
   ```
   go build app.go
   ```

2. Run the compiled binary.
   ```
   ./app
   ```

3. Server is now running. `telnet 127.0.0.1 8080` Connect the server and type your name when prompted.

## Usage

Once the chat application is up and running, you can use the following commands:  
- `/list`: This command will display a list of all online users.
- `/send <name> <message>`: Use this command to send a message to a specific user. Replace `<name>` with the recipient's name, as listed via the `/list` command, and `<message>` with your text message.
- Sending a message without specifying a command (e.g., just typing a message) will broadcast it to all online users.
- `/help`: Shows detail about chat commands.

## Concurrency

The chat application is designed to handle multiple users concurrently. It demonstrates basic concurrency features in GoLang, making it an excellent project for learning concurrent programming.

## Contributing

Feel free to use, modify, or alter the code as needed. If you'd like to contribute or suggest improvements, you are more than welcome to send a Pull Request (PR). This project is open to collaboration and learning together.

Please ensure that your contributions align with the purpose of this project, which is to serve as a learning resource for GoLang and socket programming.

Happy chatting and coding!
