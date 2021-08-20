import { AxiosResponse } from "axios";
import { useServerStatus } from "contexts/serverStatusContext";
import MaintainancePage from "pages/MaintainancePage";
import NotFound from "pages/NotFound";
import React, { Suspense } from "react";
import { Redirect, Route, Switch } from "react-router-dom";
import getStatus from "services/status";
import { Status } from "utils/enums/status";

const ShowSecret = React.lazy(() => import("pages/FetchSecret"));
const CreateSecret = React.lazy(() => import("pages/CreateSecret"));

const AppRouter: React.FC = () => {
	const [status, setServerStatus] = useServerStatus();
	const isMounted = React.useRef(true);

	React.useEffect(() => {
		getStatus().then((response: AxiosResponse) => {
			const hostname = location.hostname;
			setServerStatus({ ...response.data, host: hostname });
		});

		return () => {
			isMounted.current = false;
		};
	}, []);

	return (
		<Suspense fallback="loading...">
			<Switch>
				<Route path="/maintainance" exact component={MaintainancePage} />
				{status.status === Status.maintainance ? <Redirect to="/maintainance" /> : null}
				<Route path="/secret/:token" exact component={ShowSecret} />
				<Route path="/" exact component={CreateSecret} />
				<Route path="/not-found" component={NotFound} />
				<Redirect from="*" to="/not-found" />
			</Switch>
		</Suspense>
	);
};

export default AppRouter;
