import {React} from "../m_react.js";

export default function TextInput(props) {
    const {value, onChange, autoFocus, ...otherProps} = props;
    const callback = React.useCallback((element) => {
        if (element && autoFocus) {
            setTimeout(function () {element.focus();}, 10);
        }
    }, []);

    return <input value={value}
                  onChange={(event) => {onChange(event.target.value);}}
                  ref={callback} {...otherProps} />;
}
