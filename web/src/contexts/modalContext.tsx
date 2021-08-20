import React from "react";
import { useReducer } from "react";
import { ModalType } from "utils/enums/modal";

type State = {
	modalType?: string;
	modalProps?: Record<string, never>;
};

type Action =
	| { type: ModalType.SHOW_MODAL; payload?: Record<string, never> }
	| { type: ModalType.HIDE_MODAL; payload?: Record<string, never> };

type ModalProviderProps = {
	children: React.ReactNode;
};

const initialState: State = {
	modalType: undefined,
	modalProps: {}
};

const ModalContext = React.createContext<{ state: State; dispatch: React.Dispatch<Action> } | undefined>(undefined);

const modalReducer = (state = initialState, action: Action) => {
	switch (action.type) {
		case ModalType.SHOW_MODAL:
			return {
				...state,
				modalType: ModalType.SHOW_MODAL,
				modalProps: action.payload
			};
		case ModalType.HIDE_MODAL:
			return initialState;
		default:
			throw new Error(`unhandled action type ${(action as Action).type}`);
	}
};

const ModalProvider: React.FC<ModalProviderProps> = ({ children }) => {
	const [state, dispatch] = useReducer(modalReducer, { modalType: undefined, modalProps: {} });
	const value = { state, dispatch };

	return <ModalContext.Provider value={value}>{children}</ModalContext.Provider>;
};

const useModal = (): { state: State; dispatch: React.Dispatch<Action> } => {
	const context = React.useContext(ModalContext);

	if (!context) {
		throw new Error("useModal should be within a ModalProvider");
	}
	return context;
};

export { useModal, ModalProvider };
