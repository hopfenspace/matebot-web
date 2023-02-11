import { StylesConfig } from "react-select";

// Type used in <Select/>'s generic parameter: `Option`
export type Option = { label: string; value: number };

// Convert an api model into an option
export function apiToOption(named: { id: number; name: string }): Option {
    return { label: named.name, value: named.id };
}

const STYLES: StylesConfig<Option, false> = {
    option: (styles, { data: _1, isDisabled: _2, isFocused, isSelected: _3 }) => {
        return {
            ...styles,
            backgroundColor: isFocused ? "var(--prim)" : "var(--level-0)",
        };
    },
};

export const SELECT_PROPS = {
    className: "select",
    styles: STYLES,
};
