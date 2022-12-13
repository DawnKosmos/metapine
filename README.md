# metapine

NOT FINISHED!
metapine will enable it to run and backtest tradingstrategies in an Tradingview like enviroment with MQL5 performance.
Backtesting is working. Look in the example folder for some example
Development stopped due to FTX bankrupt


## Design Backtesting 
**backend/series/ta**
Backtesting has a modular approach and every interface often implements other interface.
This design is used to be able to backtest in different ways, adjust everything from 
trading fees to trade execution to order size etc. Without changing the main infrastructure
for new features.


### Chart, Series, Conditions Interface
**Chart** saves the OHCLV(open,high,close,low,volume which are Series) of a Chart 
**Series** save Data with numerical Values
**Conditions** save Data with Boolean Values

### Backtester Interface
Once you Setup an Enviroment the Backtester Interface is used to Add Strategies
And Filter the Trading Results.

TBD