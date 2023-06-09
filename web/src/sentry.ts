import * as Sentry from "@sentry/react";
import { BrowserTracing } from "@sentry/tracing";

const defaultTracingOrigins = ["localhost", /^\//];

const initSentry = (): void => {
    const appVersion = process.env.REACT_APP_RELEASE_VERSION;
    const gitRevision = process.env.REACT_APP_GIT_REVISION;

    // ensure environment variables app version and git revision are set
    if (!appVersion) {
        // eslint-disable-next-line no-console
        console.log("App version is not set in environment variables");
    }
    if (!gitRevision) {
        // eslint-disable-next-line no-console
        console.log("Git revision is not set in environment variables");
    }
    console.log(`AppVersion: ${appVersion || ""} - GitRevision: ${gitRevision || ""}`); // eslint-disable-line no-console

    if (process.env.REACT_APP_SENTRY_DSN) {
        let tracingOrigins = defaultTracingOrigins;
        if (process.env.REACT_APP_API_BASE_URL) {
            const origin = new URL(process.env.REACT_APP_API_BASE_URL);
            tracingOrigins = [origin.host];
        }

        const environment = process.env.REACT_APP_SENTRY_ENVIRONMENT
            ? process.env.REACT_APP_SENTRY_ENVIRONMENT
            : process.env.NODE_ENV;

        Sentry.init({
            dsn: process.env.REACT_APP_SENTRY_DSN,
            integrations: [new BrowserTracing({ tracingOrigins })],
            environment,
            tracesSampleRate: 0.25,
            release: appVersion || "",
            // mute session expired errors
            ignoreErrors: ["Session expired"],
        });

        // eslint-disable-next-line no-console
        console.log("Sentry tracing initialized");
    } else {
        // eslint-disable-next-line no-console
        console.log("no Sentry configuration available");
    }
};

export default initSentry;