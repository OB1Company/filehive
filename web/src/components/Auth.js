import React from "react";
import axios from "axios";


export const getAxiosInstance = () => {
    const csrftoken = localStorage.getItem('csrf_token');
    return axios.create({
        baseURL: "",
        headers: {"x-csrf-token": csrftoken}
    })
}