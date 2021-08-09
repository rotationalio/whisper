import Footer from "components/Footer";
import Modal from "components/Modal";
import { ModalProvider } from "contexts/modalContext";
import { BrowserRouter as Router } from "react-router-dom";
import AppRouter from "routes";
import "./App.css";

const App: React.FC = () => {
	return (
		<Router>
			<div className="App">
				<ModalProvider>
					<AppRouter />
					<Modal />
				</ModalProvider>
				<Footer />
			</div>
		</Router>
	);
};

export default App;
