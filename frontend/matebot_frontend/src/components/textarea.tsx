import React from "react";

export type TextAreaProps = {
    value: string;
    onChange: (newValue: string) => any;
    autoFocus?: boolean;
} & any;

export default function TextArea(props: TextAreaProps) {
    const { value, onChange, autoFocus, ...otherProps } = props;
    const callback = React.useCallback((element: HTMLTextAreaElement) => {
        if (element && autoFocus) {
            setTimeout(function () {
                element.focus();
            }, 10);
        } // eslint-disable-next-line
    }, []);

    return (
        <textarea
            value={value}
            onChange={(e) => {
                onChange(e.target.value);
            }}
            ref={callback}
            {...otherProps}
        />
    );
}
