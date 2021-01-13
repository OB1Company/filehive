import { createStore, compose } from "redux";

const initialState = {
    token: null,
}

export const SET_TOKEN = "SET_TOKEN";

const rootReducer = (state = initialState, action) => {
    switch (action.state) {
        case SET_TOKEN:
            return {
                token: action.payload,
            };
        default:
            return state;
    }
};

export default createStore (
    rootReducer,
    compose(
        window.__REDUX_DEVTOOLS_EXTENSION__
            ? window.__REDUX_DEVTOOLS_EXTENSION__()
            : (f) => f
    )
);