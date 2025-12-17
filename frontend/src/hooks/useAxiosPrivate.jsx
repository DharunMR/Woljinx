import { useRef, useEffect } from "react";
import axios from "axios";
import useAuth from "./useAuth";

const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () => {
  const { setAuth } = useAuth();

  // Create axios instance ONCE
  const axiosAuth = useRef(
    axios.create({
      baseURL: apiUrl,
      withCredentials: true,
    })
  ).current;

  const isRefreshing = useRef(false);
  const failedQueue = useRef([]);

  const processQueue = (error, response = null) => {
    failedQueue.current.forEach(prom => {
      error ? prom.reject(error) : prom.resolve(response);
    });
    failedQueue.current = [];
  };

  useEffect(() => {
    const responseInterceptor = axiosAuth.interceptors.response.use(
      response => response,
      async error => {
        const originalRequest = error.config;

        if (
          originalRequest?.url?.includes("/refresh") &&
          error.response?.status === 401
        ) {
          console.error("❌ Refresh token expired");
          setAuth(null);
          return Promise.reject(error);
        }

        if (
          error.response?.status === 401 &&
          !originalRequest._retry
        ) {
          if (isRefreshing.current) {
            return new Promise((resolve, reject) => {
              failedQueue.current.push({ resolve, reject });
            })
              .then(() => axiosAuth(originalRequest))
              .catch(err => Promise.reject(err));
          }

          originalRequest._retry = true;
          isRefreshing.current = true;

          try {
            await axiosAuth.post("/refresh");
            processQueue(null);
            return axiosAuth(originalRequest);
          } catch (refreshError) {
            processQueue(refreshError);
            setAuth(null);
            return Promise.reject(refreshError);
          } finally {
            isRefreshing.current = false;
          }
        }

        return Promise.reject(error);
      }
    );

    // 🧹 Cleanup interceptor
    return () => {
      axiosAuth.interceptors.response.eject(responseInterceptor);
    };
  }, [axiosAuth, setAuth]);

  return axiosAuth;
};

export default useAxiosPrivate;
