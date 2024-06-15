import React, {Component} from "react";
import './Message.scss';

class Message extends Component {
    constructor(props) {
        super((props));
        this.state = {
            message: typeof props.message === "string" ? JSON.parse(props.message) : props.message,
        };
    }

    render() {
        return(
            <div className="Message">
                {this.state.message.body}
            </div>
        );
    };
}

export default Message;