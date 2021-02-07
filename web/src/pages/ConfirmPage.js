import React from 'react'
import { Redirect, useLocation } from 'react-router-dom'
import {getAxiosInstance} from "../components/Auth";

function useQuery() {
    return new URLSearchParams(useLocation().search);
}

export default function ConfirmPage() {

    let query = useQuery();

    // Update database to confirm email
    const code = query.get("code");
    const email = query.get("email");

    if(code.length === 6) {
        const confirmEmail = async () => {
            const instance = getAxiosInstance();
            instance.get(
                "/api/v1/confirm?email="+email+"&code="+code,
            )
                .then((res)=>{
                    console.log(res);
                    return (
                        <Redirect to="/login?confirmed=1" />
                    )
                })
                .catch((err)=> {
                    console.error(err);
                });
        }
        confirmEmail();
    }

    return (
        <Redirect to="/login" />
    )

}