import React, {useState} from 'react'
import {Helmet} from "react-helmet";
import Header from "../Header";
import Footer from "../Footer";
import {Link} from "react-router-dom";
import ErrorBox, {SuccessBox} from "../components/ErrorBox";
import {getAxiosInstance} from "../components/Auth";

function validateEmail(email) {
    const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

export default function PasswordResetPage() {

    const [email, setEmail] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const HandlePasswordReset = (e)=>{
        e.preventDefault();
        setError("");

        if(!validateEmail(email)) {
            setError("Please enter a valid email address");
        }

        const submitRequest = async ()=>{
            const instance = getAxiosInstance();
            instance.get("/api/v1/passwordreset?email="+email)
                .then((response)=>{
                    console.log(response);
                    setSuccess("A password reset link has been sent to your email.");
                })
                .catch((err)=>{
                    console.log(err);
                })
        }
        submitRequest();

    }

    return (
        <div className="container">
            <Helmet>
                <title>Filehive | Password Reset</title>
            </Helmet>
            <Header/>
            <div className="subBody">
                <div className="maincontent">

                    <div className="Login form-540">
                        <h2>Password Reset</h2>

                        <form onSubmit={HandlePasswordReset}>
                            <label>
                                Email address
                                <input type="text" name="email" placeholder="Enter email" onChange={e => setEmail(e.target.value)}/>
                            </label>
                            <div>
                                <input type="submit" value="Reset password" className="raise orange-button" />
                                <Link to='/login'>Log in</Link>
                            </div>

                            {error &&
                            <ErrorBox message={error}/>
                            }
                            {success &&
                            <SuccessBox message={success}/>
                            }
                        </form>

                    </div>

                </div>
            </div>
            <Footer/>
        </div>
    )
}