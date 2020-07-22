# parametric tolerance interval
from numpy.random import seed
from numpy.random import randn
from numpy import mean
from numpy import sqrt
from scipy.stats import chi2
from scipy.stats import norm
import sys

data = [float(x) for x in sys.argv[1:]]
# specify degrees of freedom
n = len(data)
dof = n - 1
# specify data coverage
prop = 0.95 #95% of data
prop_inv = (1.0 - prop) / 2.0
gauss_critical = norm.isf(prop_inv)
# print('Gaussian critical value: %.3f (coverage=%d%%)' % (gauss_critical, prop*100))
# specify confidence
prob = 0.99 # with the confidence of 99%
chi_critical = chi2.isf(q=prob, df=dof)
# print('Chi-Squared critical value: %.3f (prob=%d%%, dof=%d)' % (chi_critical, prob*100, dof))
# tolerance
interval = sqrt((dof * (1 + (1/n)) * gauss_critical**2) / chi_critical)
# print('Tolerance Interval: %.3f' % interval)
# summarize
data_mean = mean(data)
lower, upper = data_mean-interval, data_mean+interval
# print('%.2f to %.2f covers %d%% of data with a confidence of %d%%' % (lower, upper, prop*100, prob*100))
print('%.2f,%.2f' % (lower, upper))