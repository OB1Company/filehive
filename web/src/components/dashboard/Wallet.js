import React from 'react'
import { Link } from 'react-router-dom'

export default function Wallet() {

    return (

        <div className="maincontent margins-30">
            <h2>Wallet <span className="h2-subtitle">(2.349 FIL)</span></h2>

            <div className="withdrawal-deposit-container">
                <div className="wd-container">
                    <h3>Deposit</h3>
                    <p>Send FIL to the address below to add funds to your wallet.</p>
                    <div className="qr-code-deposit"></div>
                    <p className="center"></p>
                    <Link onClick=""/>
                </div>
                <div className="wd-container">
                    <h3>Withdrawal</h3>
                    <p>Specify a FIL address below to send your funds to.</p>
                    <form>

                    </form>
                </div>
            </div>
            <div className="transaction-ledger">
                <h3>Transactions</h3>

            </div>


        </div>

    );
}