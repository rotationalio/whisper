import AppRouter from "Router";
import { BrowserRouter as Router } from "react-router-dom";
import "./App.css";

const App: React.FC = () => {
	return (
		<Router>
			<div className="App">
				Render
				<AppRouter />
			</div>
		</Router>
	);
};

export default App;
