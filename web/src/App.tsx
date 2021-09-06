import CreateSecretModal from "components/CreateSecretModal";
import { ModalProvider } from "contexts/modalContext";
import { ServerStatusProvider } from "contexts/serverStatusContext";
import { BrowserRouter as Router } from "react-router-dom";
import AppRouter from "routes";
import "./App.css";

const App: React.FC = () => {
	return (
		<Router>
			<div className="App">
				<ServerStatusProvider>
					<ModalProvider>
						<AppRouter />
						<CreateSecretModal />
					</ModalProvider>
				</ServerStatusProvider>
			</div>
		</Router>
	);
};

export default App;
